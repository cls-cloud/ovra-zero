#!/bin/bash

target=${1:-.}

if [ -d "$target" ]; then
    echo "Scanning directory: $target"
    find "$target" -type f -path "*/model/*" -name "*gen.go" | while read -r file; do
        bash "$0" "$file"
    done
    echo "Cleaning up backup files (*.bak)..."
    find "$target" -type f -name "*.bak" -delete
    echo "All replacements done. Backups removed."
    exit 0
fi

file="$target"

if grep -qE '^[[:space:]]*DelFlag[[:space:]]+[a-zA-Z0-9_]+' "$file"; then
    echo "🛠  Processing: $file"

    sed -i.bak -E \
        's/^[[:space:]]*DelFlag[[:space:]]+[a-zA-Z0-9_]+[[:space:]]+`[^`]*`.*$/\tDelFlag soft_delete.DeletedAt `gorm:"softDelete:flag"`/' \
        "$file"

    if ! grep -q 'gorm.io/plugin/soft_delete' "$file"; then
        echo "➕  Adding import for soft_delete in $file"
        awk '
            /^import[[:space:]]*\(/ && !added {
                print
                print "\t\"gorm.io/plugin/soft_delete\""
                added=1
                next
            }
            { print }
        ' "$file" > "${file}.tmp" && mv "${file}.tmp" "$file"
    fi
fi
