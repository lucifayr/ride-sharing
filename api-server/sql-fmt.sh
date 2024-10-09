for file in $(find db -type f -iname '*.sql'); do
    sql-formatter --config .sql-formatter.json --fix "$file"
done
