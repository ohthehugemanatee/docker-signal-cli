# Signal Nixplay bridge

This project creates a bridge between my Signal messenger Groups and my Nixplay account. It simply logs into my signal, and when it sees an image (JPEG, > 512k) in the target group, it automatically emails it to the email address of your choice - in my case, my Nixplay address. It's basically a mail wrapper around the excellent [signal-cli](https://github.com/AsamK/signal-cli).

**Obvious note: Files sent over email lose the privacy protetions of Signal!** Be a good person and warn any others in the chat that their images are going through your (or google's, or microsoft's...) SMTP server to your photo frame.

## ToDo
* Make first time setup easier


## Installation/Usage
* Copy the `env.example` file to `env`, and modify it with your own configuration. If you don't know your group ID yet, make one up while you go through the registration process.
```
$> cp env.example .env
$> nano .env
```
* Register yourself either as a new phone number (make sure to add it to the Group!), or as a "Slave device" on an existing phone number.
** New number: `run.sh signal-cli -u <your phone number> register`, and then later `run.sh signal-cli -u <your phone number> validate <your code from SMS>`. You can register with `--voice` if the phone can't receive SMS.
** Slave device: `run.sh signal-cli link -n "optional device name"`. Pass the resulting `tsdevice:/` link through a QR code generator to scan it with your phone.
* Once you are registered, look up the groups you belong to with `run.sh signal-cli -u <your phone number> listGroups`. Find the ID of the group you want to monitor, and add it to your `.env` file.
* Run the service with `run.sh`


Note: Your phone number must always include the country code. e.g. `+49123456789`. In North America it's `+1`.
