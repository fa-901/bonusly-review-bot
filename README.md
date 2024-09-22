# bonusly-review-bot
Simple bot that automatically rewards your colleagues Bonusly points for reviewing your pull requests.

### How does it work?
The bot scans your open pull requests, identifies recent reviewers, and automatically sends them Bonusly points on your behalf.

## Setup
- You will need to generate 2 Access tokens and set them in your environment variables.

  - **GitHub Personal Access Token**: Create a GitHub PAT from [here](https://github.com/settings/tokens). Make sure to grant **all** `read` permissions to access complete pull request information. _After token is created, click on 'Conifgure SSO' and authorize your organization._ 
  - **Bonusly Access Token**: Create a Bonusly access token from [here](https://bonus.ly/api_keys/new).

- Run `main.go`

## Troubleshooting
- `I have reviewed someone's PR, but didn't get any points`
  - Your GitHub email and Bonusly email is not the same.
  - Your email address is private. The bot will not be able to find your email in order to find your Bonusly account.
    - If your email address is private, the bot will attempt to find your email from any commits pushed with your GitHub account.

- `My email is private and I have no commits, what now?`
  - In that case, the bot will attempt to find your email using Bonusly autocomplete request.
  - **Note**: This method is unreliable and your Bonusly points may be awarded to the wrong person with a similar name.

- `I do not have a name in my GitHub account / does not match with Bonusly.`
  - Nothing to do here 🤷. If you have a solution, feel free to open a PR though.

- Add more FAQs later