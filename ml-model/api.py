from flask import Flask, request, jsonify
from credit_scoring_model import CreditScoringModel
import pandas as pd

app = Flask(__name__)
model = CreditScoringModel()
model.load_model('credit_scoring_model.joblib')

@app.route('/predict', methods=['POST'])
def predict():
    try:
        data = request.json
        features = pd.DataFrame([data])
        prediction = model.predict(features)[0]
        
        if prediction >= 700:
            risk_category = 'Low'
        elif prediction >= 640:
            risk_category = 'Medium'
        else:
            risk_category = 'High'
        
        return jsonify({
            'credit_score': float(prediction),
            'risk_category': risk_category
        }), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 400

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)