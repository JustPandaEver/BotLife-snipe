<-- Guide -->

1. Open config.yml (best editor to use is notepad++)

    a. Host can be left unchanged unless you have a node you would like to use.

    b. Contract, bought, pcsaddress, and bnbaddress should be left unchanged.

    c. Enter your private key in the private key field. For maximum security you should create a new wallet then transfer 
        only the funds you need and use this wallet.

    d. Enter the amount in terms of BNB you would like to purchase in amountin.

    e. Adjust gaslimit and gasprice to your liking but they should be fine as they are.

    f. Enter your Bot API key found in the next section.

2. Setup telegram

    a. Search botfather on telegram

    b. Send the message /newbot and choose a name for your bot. This can be anything, then choose a username, it can be anything but
        must end in bot.

    c. Copy and paste the access token into the config.yml file in tgbotapi

    d. Send the message "/setprivacy" to BotFather, then send "@botusername", then "disable".

    e. Create a new group, naming it anything. Add your bot by searching @botusername 

    f. Search telefeed then send the message /start to the bot. Click connect and follow the guide to link your account.

    g. Send the message "/chats 'phoneNumber'" where phoneNumber is the phone number you registered with. Do not include the quotes
    
    h. Copy and save somewhere the numbers beside the group you created with your bot AND the channel or group you would like to watch.

    i. Send the message "/redirection add group1 on 'phoneNumber'" where phoneNumber is the phone number you registered with. Do not include the quotes

    j. Send the message "'source' - 'destination'" where source is the id of the target group and destination is the id of the group you created.
        Do not include quotes.

    k. Make sure to save config.yml

3. Run bot by either double clicking or navigating to the folder then executing the bot.

    a. Enter your license key, make sure there are no extra white spaces.

4. Run the bot again. You should get a message saying starting...

If you have any questions or issues email support@wisdom-bots.com or join the telegram t.me/wisdombotsgroup