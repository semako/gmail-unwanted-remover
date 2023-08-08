# Gmail Unwanted Remover

Purpose: on regular basis permanently delete unwanted email, like: spam, trash, etc.

I personally follow rule: empty inbox, so software reading `inbox` and `spam` folders. So my inbox contains a small amount of emails usually,
so sorry, no pagination.

- Trying to find emails in `inbox` where:
   - substring occurs
   - from contains specified domain
   - from contains exact email
   - permanently delete them 
- Permanently delete all in `spam` folder

## Prepare Google OAuth
- Follow [instructions](https://developers.google.com/gmail/api/quickstart/go#configure_the_oauth_consent_screen) on official page to create `credentials.json` file
- Put `credentials.json` to root folder

## Setup
- Copy `cmd/config.dist.yml` to `cmd/config.yml`
- Make necessary changes there
- Run `go build -o gmail-unwanted-remover cmd/gmail-unwanted-remover/main.go`
- First run: `./gmail-unwanted-remover`, it will show auth url, follow it
- Copy auth code from url in browser
- Put it to console, press Enter (now you have `token.json` -- to make authorized requests to Gmail API)
  - DO NOT STORE `credentials.json` and `token.json` in git repository

## Run as Daemon on Linux systems (example)
- Install supervisord `sudo apt install supervisor`
- Create a config file: `/etc/supervisor/conf.d/gmail-unwanted-remover` with contents
  ```ini
  [program:gmail-unwanted-remover]
  command=<path to binary>
  directory=<dir where binary stored>
  autostart=true
  autorestart=true
  startretries=3
  numprocs=1
  startsecs=0
  process_name=%(program_name)s_%(process_num)02d
  stderr_logfile=/var/log/supervisor/%(program_name)s_stderr.log
  stderr_logfile_maxbytes=10MB
  stdout_logfile=/var/log/supervisor/%(program_name)s_stdout.log
  stdout_logfile_maxbytes=10MB
  ```
- Run `sudo supervisorctl update` to apply changes
- Check if it is running: `sudo supervisorctl status`
  ```bash
  gmail-unwanted-remover:gmail-unwanted-remover_00         RUNNING   pid 110045, uptime 0:01:55
  ```
