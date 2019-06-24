UPDATE customers
SET company_name = $2, company_address = $3, tax_id = $4, email = $5, phone_number = $6
WHERE id = $1;