import csv
import os
import random
import datetime
import requests
import faker

NUM_RECORDS = 5000
NUM_USERS = 200

input_csv_headers = ['id', 'user_id', 'amount', 'datetime']
output_csv_headers = ['user_id', 'balance', 'total_debts', 'total_credits']
count_csv_headers = ['total_transactions']
USER_CREATION_URL = 'http://127.0.0.1:8000/user-balance-api/users/create'

def ensure_files(input_file, output_file, count_file):
    os.makedirs(os.path.dirname(count_file), exist_ok=True)

    if os.path.exists(input_file):
        os.remove(input_file)
    if os.path.exists(output_file):
        os.remove(output_file)
    if not os.path.exists(count_file):
        with open(count_file, 'w', newline='') as f:
            writer = csv.writer(f)
            writer.writerow(count_csv_headers)
            writer.writerow([0])  # Initialize with 0 transactions if file doesn't exist

def create_user():
    fake = faker.Faker()
    first_name = fake.first_name()
    last_name = fake.last_name()
    email = f"{first_name.lower()}.{last_name.lower()}@example.com"

    user_data = {
        "first_name": first_name,
        "last_name": last_name,
        "email": email
    }

    response = requests.post(USER_CREATION_URL, json=user_data)

    if response.status_code == 201:
        print(f"User {first_name} {last_name} created successfully.")
        return response.json().get('user_id')
    else:
        print(f"Failed to create user {first_name} {last_name}. Status code: {response.status_code}")
        return None

def create_users(num_users):
    ids = []

    for _ in range(num_users):
        user_id = create_user()
        if user_id:
            ids.append(user_id)

    return ids

def calculate_total_transactions(count_file):
    total_transactions = 0

    if os.path.exists(count_file):
        with open(count_file, 'r') as f:
            reader = csv.reader(f)
            next(reader)  # Skip header
            total_transactions = int(next(reader)[0])  # Read the total transaction count

    return total_transactions

def update_transaction_count(count_file, new_transactions):
    total_transactions = calculate_total_transactions(count_file) + new_transactions

    with open(count_file, 'w', newline='') as f:
        writer = csv.writer(f)
        writer.writerow(count_csv_headers)
        writer.writerow([total_transactions])

def generate_input_csv(input_file, user_ids, start_transaction_id):
    records = []

    for i in range(start_transaction_id, start_transaction_id + NUM_RECORDS):
        user_id = random.choice(user_ids)
        amount = round(random.uniform(-1000, 1000), 2)  # random amount between -1000 and 1000
        random_date = datetime.datetime(2023, 1, 1) + datetime.timedelta(days=random.randint(0, 548))
        datetime_str = random_date.strftime('%Y-%m-%dT%H:%M:%SZ')

        record = [i, user_id, amount, datetime_str]
        records.append(record)

    with open(input_file, 'w', newline='') as f:
        writer = csv.writer(f)
        writer.writerow(input_csv_headers)
        writer.writerows(records)

def calculate_output(input_file, output_file):
    user_data = {}

    with open(input_file, 'r') as f:
        reader = csv.DictReader(f)
        for row in reader:
            user_id = row['user_id']
            amount = float(row['amount'])

            if user_id not in user_data:
                user_data[user_id] = {'balance': 0, 'total_debts': 0, 'total_credits': 0}

            user_data[user_id]['balance'] += amount
            if amount < 0:
                user_data[user_id]['total_debts'] += 1
            else:
                user_data[user_id]['total_credits'] += 1

    sorted_user_data = sorted(user_data.items(), key=lambda x: int(x[0]))

    with open(output_file, 'w', newline='') as f:
        writer = csv.writer(f)
        writer.writerow(output_csv_headers)
        for user_id, data in sorted_user_data:
            writer.writerow([user_id, round(data['balance'], 2), round(data['total_debts'], 2), round(data['total_credits'], 2)])

if __name__ == "__main__":
    file_path = 'scripts/generate_transactions/files/'
    input_file = f'{file_path}input_data.csv'
    output_file = f'{file_path}expected_output_data.csv'
    output_count_file = f'{file_path}transactions_count.csv'

    ensure_files(input_file, output_file, output_count_file)
    total_transactions_before = calculate_total_transactions(output_count_file)

    user_ids = create_users(NUM_USERS)
    if len(user_ids) < NUM_USERS:
        print("Some users could not be created. Exiting...")
        exit(1)

    generate_input_csv(input_file, user_ids, total_transactions_before + 1)
    calculate_output(input_file, output_file)
    update_transaction_count(output_count_file, NUM_RECORDS)

    print(f"Input CSV generated at: {input_file}")
    print(f"Expected output CSV generated at: {output_file}")
    print(f"Transaction count CSV updated at: {output_count_file}")