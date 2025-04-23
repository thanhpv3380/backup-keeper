# Backup Keeper

Backup Keeper is a tool designed to automate the process of collecting data from a MongoDB database, compressing it, and uploading it to Google Drive. It also supports notifications via Telegram to keep you informed about the backup process.

## Features

- Collect data from MongoDB collections.
- Compress data into `.zip` files.
- Upload backups to Google Drive.
- Schedule backups using cron expressions.
- Send notifications via Telegram.

## Configuration

The project uses a `.env` file to manage configuration. Ensure you create a `.env` file in the root directory with the following structure:

```env
# MongoDB Configuration
MONGODB_URI=mongodb://<username>:<password>@<host>:<port>/<database>
MONGODB_DATABASE=<database_name>

# Telegram Configuration
TELEGRAM_BOT_TOKEN=<your_telegram_bot_token>
TELEGRAM_CHAT_ID=<your_telegram_chat_id>

# Google Drive Configuration
GOOGLE_DRIVE_CREDENTIALS_FILE=<google_service_account_credential>
GOOGLE_DRIVE_FOLDER_ID=<google_drive_folder_id>

# Backup Configuration
BACKUP_DATA_SOURCE=<data_source_name> # Show data source in Telegram message
BACKUP_CRON_SCHEDULE=0 0 * * * # Example: Run daily at midnight
BACKUP_TIMEZONE=Asia/Ho_Chi_Minh
```

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/thanhpv3380/backup-keeper.git
   cd backup-keeper
   ```

2. **Install dependencies**:
   Ensure you have Go installed (version 1.18 or later). Run the following command to install dependencies:

   ```bash
   go mod tidy
   ```

3. **Set up the `.env` file**:
   Create a `.env` file in the root directory and configure it as shown in the [Configuration](#configuration) section.

4. **Set up Google Drive credentials**:
   - Download your Google Drive service account credentials JSON file.
   - Place it in the `credentials` directory (e.g., `credentials/drive_credentials.json`).

---

## Running the Application

To run the application, use the following command:

```bash
go run ./cmd/main.go
```

### Running with Cron Scheduling

If you want to schedule backups automatically based on the cron expression defined in the `.env` file, ensure the `BACKUP_CRON_SCHEDULE` is properly set. The application will handle the scheduling internally.

---

## Example Output

When the application runs successfully, you will see logs like the following:

```plaintext
2025/04/23 14:00:00 Starting backup process...
2025/04/23 14:00:01 MongoDB collector initialized and connection verified
2025/04/23 14:00:02 Processing collection: users
2025/04/23 14:00:03 Uploaded backup_20250423_140003.zip to Google Drive folder xxx
2025/04/23 14:00:04 ✅ Backup successful: backup_20250423_140003.zip
```

Let me know if you need further assistance!
