DROP TABLE IF EXISTS loan_investments;
DROP TABLE IF EXISTS loan_approvals;
DROP TABLE IF EXISTS loan_disbursements;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS accounts;

CREATE TABLE accounts (
	id BIGSERIAL PRIMARY KEY,
	customer_xid VARCHAR ( 36 ) UNIQUE NOT NULL,
	name VARCHAR ( 36 ) NOT NULL,
	email VARCHAR ( 36 ) UNIQUE NOT NULL,
	password VARCHAR ( 36 ) NOT NULL
);

CREATE TABLE employees (
	id BIGSERIAL PRIMARY KEY,
	employee_xid VARCHAR ( 36 ) UNIQUE NOT NULL,
	name VARCHAR ( 36 ) NOT NULL,
	email VARCHAR ( 36 ) UNIQUE NOT NULL,
	password VARCHAR ( 36 ) NOT NULL
);

CREATE TABLE loans (
	id BIGSERIAL PRIMARY KEY,
	loan_xid VARCHAR ( 36 ) UNIQUE NOT NULL,
	borrower_id BIGINT NOT NULL,
	principal_amount INT NOT NULL,
	rate INT NOT NULL,
	roi INT NOT NULL,
	status VARCHAR ( 10 ) NOT NULL,
	FOREIGN KEY (borrower_id) REFERENCES accounts(id)
);

CREATE TABLE loan_approvals (
	id BIGSERIAL PRIMARY KEY,
	loan_id BIGINT NOT NULL,
	approved_by BIGINT NOT NULL,
	approved_date DATE NOT NULL,
	approval_proof_image_path TEXT NOT NULL,
	FOREIGN KEY (approved_by) REFERENCES employees(id)
);

CREATE TABLE loan_investments (
	id BIGSERIAL PRIMARY KEY,
	loan_id BIGINT NOT NULL,
	amount INT NOT NULL,
	invested_by BIGINT NOT NULL,
	agreement_pdf_path TEXT NOT NULL,
	FOREIGN KEY (invested_by) REFERENCES accounts(id)
);

CREATE TABLE loan_disbursements (
	id BIGSERIAL PRIMARY KEY,
	loan_id BIGINT NOT NULL,
	employee_id BIGINT NOT NULL,
	agreement_pdf_path TEXT NOT NULL,
	disbursed_date DATE NOT NULL,
	FOREIGN KEY (employee_id) REFERENCES employees(id)
);