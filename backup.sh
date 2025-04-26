#!/bin/bash

# **获取脚本所在目录**
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)

# **备份目录**
BACKUP_BASE="$SCRIPT_DIR/appdata/backup"
DB_SOURCE="$SCRIPT_DIR/appdata/db"
# **日志文件路径**
LOG_FILE="$SCRIPT_DIR/appdata/db_backup.log"
# # 部署
# BACKUP_BASE="$SCRIPT_DIR/backupdb"
# DB_SOURCE="$SCRIPT_DIR/ownsadb"
# # 部署日志文件路径
# LOG_FILE="$SCRIPT_DIR/ownsadb/db_backup.log"
# **获取当前年月**
CURRENT_YEAR=$(date +'%Y')
CURRENT_MONTH=$(date +'%m')
# **计算上个月**
if [ "$CURRENT_MONTH" -eq 1 ]; then
    LAST_YEAR=$((CURRENT_YEAR - 1))
    LAST_MONTH=12
else
    LAST_YEAR=$CURRENT_YEAR
    LAST_MONTH=$((CURRENT_MONTH - 1))
fi


DATE="${LAST_YEAR}-${LAST_MONTH#0}"
BACKUP_DIR="$BACKUP_BASE/$DATE"

# **创建备份目录**
mkdir -p "$BACKUP_DIR"

# **备份数据库**
cp -r "$DB_SOURCE"/* "$BACKUP_DIR"

# **记录日志**
echo "$(date +'%Y-%m-%d %H:%M:%S') - 备份完成，存储在：$BACKUP_DIR" >> "$LOG_FILE"

# **清理 2 年前的旧备份**
find "$BACKUP_BASE" -type d -mtime +365 -exec rm -rf {} \;

echo "✅ 备份完成，存储在：$BACKUP_DIR"
