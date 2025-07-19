#!/bin/bash
input_file=$1
temp_file=$(mktemp)

# Фильтруем coverage файл
grep -v -E "\.pb\.go:|queries\.sql\.go:|_mock\.go:|internal/server/repositories/database/generated/|gophkeeper/internal/protos|gophkeeper/internal/errs|gophkeeper/cmd" "$input_file" > "$temp_file"

# Показываем результат по пакетам
go tool cover -func="$temp_file" | grep -E "^[a-zA-Z].*total:" | sort

# Общее покрытие
echo "---"
go tool cover -func="$temp_file" | tail -1

rm "$temp_file"