import numpy as np
import pandas as pd

def generate_synthetic_credit_data(n_samples=5000):
    np.random.seed(42)
    
    data = {
        'income_level': np.random.lognormal(10, 1, n_samples),
        'debt_level': np.random.uniform(0, 100000, n_samples),
        'credit_utilization': np.random.beta(2, 5, n_samples),
        'credit_history_length': np.random.randint(0, 360, n_samples),
        'num_credit_accounts': np.random.randint(0, 10, n_samples),
        'num_credit_inquiries': np.random.poisson(2, n_samples),
        'age': np.random.randint(18, 80, n_samples),
        'payment_history': np.random.choice(['Good', 'Fair', 'Poor'], n_samples, p=[0.7, 0.2, 0.1]),
        'employment_status': np.random.choice(['Employed', 'Self-Employed', 'Unemployed'], n_samples, p=[0.8, 0.15, 0.05]),
        'education_level': np.random.choice(['High School', 'Bachelor', 'Master', 'PhD'], n_samples, p=[0.3, 0.4, 0.2, 0.1])
    }
    
    df = pd.DataFrame(data)
    
    # Create target variable with some logical rules
    df['credit_score'] = (
        df['income_level'] / 10000 * 30 +
        (1 - df['credit_utilization']) * 300 +
        df['credit_history_length'] / 5 +
        (df['payment_history'] == 'Good') * 100 -
        (df['payment_history'] == 'Poor') * 100 -
        df['num_credit_inquiries'] * 5 +
        (df['employment_status'] == 'Employed') * 50
    ).clip(300, 850).astype(int)
    
    return df

# Generate synthetic data
synthetic_data = generate_synthetic_credit_data(10000)
print(synthetic_data.head())
print(synthetic_data.describe())

# Save to CSV for later use
synthetic_data.to_csv('synthetic_credit_data.csv', index=False)