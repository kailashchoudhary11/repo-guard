import os
import re

import numpy as np
import yaml
from flask import Flask, jsonify, request
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

app = Flask(__name__)

# BAAI/bge-m3 for state-of-the-art accuracy
model = SentenceTransformer("BAAI/bge-m3", device="cpu")

TITLE_WEIGHT = 0.7
BODY_WEIGHT = 0.3


def clean(text: str) -> str:
    text = re.sub(r"```.*?```", "", text, flags=re.S)
    text = re.sub(r"\s+", " ", text)
    return text.strip()


def issue_text(title: str, body: str) -> str:
    title_clean = clean(title) if title else ""
    body_clean = clean(body) if body else ""

    if not body_clean or body_clean == "...":
        return f"GitHub Issue: {title_clean}"

    return f"GitHub Issue: {title_clean}. {body_clean}"


def parse_md_template(content: str) -> set[str]:
    """Extract boilerplate lines from a markdown issue template."""
    lines = content.split("\n")
    boilerplate = set()

    # Strip YAML front matter
    in_frontmatter = False
    body_lines = []
    for line in lines:
        stripped = line.strip()
        if stripped == "---":
            in_frontmatter = not in_frontmatter
            continue
        if in_frontmatter:
            continue
        body_lines.append(stripped)

    for line in body_lines:
        if not line:
            continue
        boilerplate.add(line)

    return boilerplate


def parse_yml_template(content: str) -> set[str]:
    """Extract boilerplate lines from a YAML form issue template."""
    boilerplate = set()
    try:
        data = yaml.safe_load(content)
    except yaml.YAMLError:
        return boilerplate

    if not isinstance(data, dict):
        return boilerplate

    body = data.get("body", [])
    if not isinstance(body, list):
        return boilerplate

    for item in body:
        if not isinstance(item, dict):
            continue

        # Labels become ### headings in rendered issues
        attrs = item.get("attributes", {})
        if isinstance(attrs, dict):
            label = attrs.get("label", "")
            if label:
                boilerplate.add(f"### {label}")

            desc = attrs.get("description", "")
            if desc:
                for line in desc.split("\n"):
                    stripped = line.strip()
                    if stripped:
                        boilerplate.add(stripped)

            placeholder = attrs.get("placeholder", "")
            if placeholder:
                for line in placeholder.split("\n"):
                    stripped = line.strip()
                    if stripped:
                        boilerplate.add(stripped)

        # Markdown blocks contain static content
        if item.get("type") == "markdown":
            value = attrs.get("value", "") if isinstance(attrs, dict) else ""
            if value:
                for line in value.split("\n"):
                    stripped = line.strip()
                    if stripped:
                        boilerplate.add(stripped)

    return boilerplate


def extract_boilerplate(templates: list[dict]) -> set[str]:
    """Parse all templates and return combined set of boilerplate lines."""
    boilerplate = set()
    for tmpl in templates:
        name = tmpl.get("name", "")
        content = tmpl.get("content", "")
        if not content:
            continue

        if name.endswith(".yml") or name.endswith(".yaml"):
            boilerplate |= parse_yml_template(content)
        elif name.endswith(".md"):
            boilerplate |= parse_md_template(content)

    return boilerplate


def strip_boilerplate(body: str, boilerplate: set[str]) -> str:
    """Remove template boilerplate lines from an issue body."""
    if not boilerplate or not body:
        return body

    cleaned_lines = []
    for line in body.split("\n"):
        stripped = line.strip()
        if stripped and stripped not in boilerplate:
            cleaned_lines.append(line)

    return "\n".join(cleaned_lines)


@app.route("/compare_issues", methods=["POST"])
def compare_issues():
    data = request.get_json(force=True)
    issue1_text = issue_text(
        data.get("issue1_title", ""),
        data.get("issue1_body", ""),
    )
    issue2_text = issue_text(
        data.get("issue2_title", ""),
        data.get("issue2_body", ""),
    )

    # BGE-M3 supports passage encoding
    embeddings = model.encode(
        [issue1_text, issue2_text],
        convert_to_numpy=True,
        normalize_embeddings=True,
        show_progress_bar=False,
    )

    similarity = float(
        cosine_similarity(
            embeddings[0].reshape(1, -1),
            embeddings[1].reshape(1, -1),
        )[0][0]
    )

    # Adjusted thresholds based on real-world GitHub issue comparison
    # These should be tuned based on your specific dataset
    label = (
        "duplicate"
        if similarity >= 0.75
        else "possible"
        if similarity >= 0.65
        else "different"
    )

    return jsonify(
        {
            "similarity": similarity,
            "label": label,
        }
    ), 200


@app.route("/batch_compare", methods=["POST"])
def batch_compare():
    data = request.get_json(force=True)

    current = data["current_issue"]
    others = data.get("other_issues", [])
    threshold = data.get("threshold", 0.85)
    templates = data.get("templates", [])

    if not others:
        return jsonify({"similar_issues": []}), 200

    # Extract boilerplate from repo's issue templates
    boilerplate = extract_boilerplate(templates)

    # Build separate title and body lists for weighted encoding
    all_issues = [current] + others
    titles = []
    bodies = []
    for issue in all_issues:
        title = clean(issue.get("title", "")) if issue.get("title") else ""
        titles.append(title if title else "untitled")

        body = issue.get("body", "")
        if boilerplate:
            body = strip_boilerplate(body, boilerplate)
        body = clean(body) if body else ""
        bodies.append(body)

    # Encode all titles in one batch
    title_embeddings = model.encode(
        titles,
        convert_to_numpy=True,
        normalize_embeddings=True,
        show_progress_bar=False,
        batch_size=64,
    )

    # Encode all cleaned bodies in one batch
    body_embeddings = model.encode(
        bodies,
        convert_to_numpy=True,
        normalize_embeddings=True,
        show_progress_bar=False,
        batch_size=64,
    )

    # Weighted combination: title-heavy when body exists, title-only when body is empty
    combined = np.zeros_like(title_embeddings)
    for i in range(len(all_issues)):
        if bodies[i] and bodies[i] != "...":
            combined[i] = TITLE_WEIGHT * title_embeddings[i] + BODY_WEIGHT * body_embeddings[i]
        else:
            combined[i] = title_embeddings[i]

    # Re-normalize combined embeddings
    norms = np.linalg.norm(combined, axis=1, keepdims=True)
    norms[norms == 0] = 1  # avoid division by zero
    combined = combined / norms

    # Cosine similarity of current (index 0) against all others
    current_embedding = combined[0].reshape(1, -1)
    other_embeddings = combined[1:]
    similarities = cosine_similarity(current_embedding, other_embeddings)[0]

    # Filter by threshold
    similar_issues = []
    for i, sim in enumerate(similarities):
        if sim >= threshold:
            similar_issues.append({
                "number": others[i]["number"],
                "similarity": float(sim),
            })

    similar_issues.sort(key=lambda x: x["similarity"], reverse=True)

    return jsonify({"similar_issues": similar_issues}), 200


if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8333))
    app.run(host="0.0.0.0", port=port)
