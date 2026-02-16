import os
import re

from flask import Flask, jsonify, request
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

app = Flask(__name__)

# BAAI/bge-m3 for state-of-the-art accuracy
model = SentenceTransformer("BAAI/bge-m3", device="cpu")


def clean(text: str) -> str:
    text = re.sub(r"```.*?```", "", text, flags=re.S)
    text = re.sub(r"\s+", " ", text)
    return text.strip()


def issue_text(title: str, body: str) -> str:
    # Combine title and body, handling empty body case
    title_clean = clean(title) if title else ""
    body_clean = clean(body) if body else ""

    # If body is empty or minimal, rely more on title
    if not body_clean or body_clean == "...":
        return f"GitHub Issue: {title_clean}"

    return f"GitHub Issue: {title_clean}. {body_clean}"


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

    if not others:
        return jsonify({"similar_issues": []}), 200

    # Build text list: current issue first, then all others
    texts = [issue_text(current.get("title", ""), current.get("body", ""))]
    for issue in others:
        texts.append(issue_text(issue.get("title", ""), issue.get("body", "")))

    # Single batch encode — all texts processed in one call
    embeddings = model.encode(
        texts,
        convert_to_numpy=True,
        normalize_embeddings=True,
        show_progress_bar=False,
        batch_size=64,
    )

    # Cosine similarity of current (index 0) against all others
    current_embedding = embeddings[0].reshape(1, -1)
    other_embeddings = embeddings[1:]
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
