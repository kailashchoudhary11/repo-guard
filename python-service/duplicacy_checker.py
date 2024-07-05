from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer, util

app = Flask(__name__)
model = SentenceTransformer('paraphrase-MiniLM-L6-v2')

@app.route('/compare_issues', methods=['POST'])
def compare_issues():
    data = request.json
    print("The data is", data)
    issue1_title = data.get('issue1_title')
    issue1_body = data.get('issue1_body')
    issue2_title = data.get('issue2_title')
    issue2_body = data.get('issue2_body')

    issue1_text = issue1_title + " " + issue1_body
    issue2_text = issue2_title + " " + issue2_body

    embedding1 = model.encode(issue1_text, convert_to_tensor=True)
    embedding2 = model.encode(issue2_text, convert_to_tensor=True)
    cosine_similarity = util.pytorch_cos_sim(embedding1, embedding2)

    return jsonify({"similarity": cosine_similarity.item()}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
