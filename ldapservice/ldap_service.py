import os
import logging
from ldap3 import Server, Connection, ALL
import psycopg2
from psycopg2.extras import execute_values

logging.basicConfig(level=logging.INFO, format='[%(levelname)s] %(message)s')

LDAP_SERVER = os.environ.get('LDAP_SERVER', 'ldap://localhost')
LDAP_USER = os.environ.get('LDAP_USER')
LDAP_PASSWORD = os.environ.get('LDAP_PASSWORD')
LDAP_BASE_DN = os.environ.get('LDAP_BASE_DN', '')

PGHOST = os.environ.get('PGHOST', 'localhost')
PGPORT = os.environ.get('PGPORT', '5432')
PGDATABASE = os.environ.get('PGDATABASE', 'agendador')
PGUSER = os.environ.get('PGUSER', 'postgres')
PGPASSWORD = os.environ.get('PGPASSWORD', '')

USER_FILTER = os.environ.get('LDAP_USER_FILTER', '(objectClass=person)')


def fetch_users():
    server = Server(LDAP_SERVER, get_info=ALL)
    logging.info('Connecting to LDAP %s', LDAP_SERVER)
    with Connection(server, user=LDAP_USER, password=LDAP_PASSWORD, auto_bind=True) as conn:
        logging.info('Fetching users with filter %s', USER_FILTER)
        conn.search(LDAP_BASE_DN, USER_FILTER, attributes=['cn', 'givenName', 'sn', 'mail'])
        users = []
        for entry in conn.entries:
            user = {
                'username': str(entry.cn),
                'first_name': str(entry.givenName or ''),
                'last_name': str(entry.sn or ''),
                'email': str(entry.mail or ''),
            }
            users.append(user)
        logging.info('Fetched %d users', len(users))
        return users


def sync_users(users):
    logging.info('Connecting to PostgreSQL %s', PGHOST)
    conn = psycopg2.connect(
        host=PGHOST,
        port=PGPORT,
        dbname=PGDATABASE,
        user=PGUSER,
        password=PGPASSWORD,
    )
    cur = conn.cursor()
    cur.execute(
        """
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT UNIQUE,
            first_name TEXT,
            last_name TEXT,
            email TEXT
        )
        """
    )
    records = [(u['username'], u['first_name'], u['last_name'], u['email']) for u in users]
    query = "INSERT INTO users (username, first_name, last_name, email) VALUES %s ON CONFLICT (username) DO UPDATE SET first_name=EXCLUDED.first_name, last_name=EXCLUDED.last_name, email=EXCLUDED.email"
    execute_values(cur, query, records)
    conn.commit()
    cur.close()
    conn.close()
    logging.info('Synchronized %d users', len(users))


def main():
    users = fetch_users()
    if users:
        sync_users(users)


if __name__ == '__main__':
    main()
