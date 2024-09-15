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
USER_CREATION_URL = 'http://127.0.0.1:8000/user-balance-api/users/create'

def ensure_files(input_file, output_file):
    if os.path.exists(input_file):
        os.remove(input_file)
    if os.path.exists(output_file):
        os.remove(output_file)

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

def generate_input_csv(input_file, user_ids):
    users = [str(i) for i in range(1, NUM_USERS + 1)]
    records = []

    for i in range(1, NUM_RECORDS + 1):
        user_id = random.choice(user_ids)
        amount = round(random.uniform(-1000, 1000), 2)  # random amount between -1000 and 1000
        random_date = datetime.datetime(2024, 1, 1) + datetime.timedelta(days=random.randint(0, 548))
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
    file_path = 'files/'
    input_file = f'{file_path}input_data.csv'
    output_file = f'{file_path}expected_output_data.csv'

    user_ids = create_users(NUM_USERS)
    if len(user_ids) < NUM_USERS:
        print("Some users could not be created. Exiting...")
        exit(1)
    ensure_files(input_file, output_file)
    generate_input_csv(input_file, user_ids)
    calculate_output(input_file, output_file)

    print(f"Input CSV generated at: {input_file}")
    print(f"Expected output CSV generated at: {output_file}")