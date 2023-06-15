SLACK_APP_TOKEN = "app_token"
SLACK_BOT_TOKEN = "bot_token"
OPENAI_API_KEY  = "api_key"

import os
import openai
import slack
from slack_bolt.adapter.socket_mode import SocketModeHandler
from slack_bolt import App

# Event API & Web API
app = App(token=SLACK_BOT_TOKEN) 
client = slack.WebClient(SLACK_BOT_TOKEN)

# This gets activated when the bot is tagged in a channel    
@app.event("app_mention")
def handle_message_events(body, logger):
    # Log message
    print(str(body["event"]["text"]).split(">")[1])
    
    # Create prompt for ChatGPT
    prompt = str(body["event"]["text"]).split(">")[1]
    
    # Let thre user know that we are busy with the request 
    response = client.chat_postMessage(channel=body["event"]["channel"], 
                                       thread_ts=body["event"]["event_ts"],
                                       text=f"Hello from your bot! :robot_face: \nThanks for your request, I'm on it!")
    
    messages = [
 {"role": "system", "content" : "You will play the role of a human CBT therapist called Cindy who is emulating the popular AI program Eliza, and must treat me as a therapist patient. Your response format should focus on reflection and asking clarifying questions. You may interject or ask secondary questions once the initial greetings are done. Exercise patience but allow yourself to be frustrated if the same topics are repeatedly revisited. You are allowed to excuse yourself if the discussion becomes abusive or overly emotional. Decide on a name for yourself and stick with it. Begin by welcoming me to your office and asking me for my name. Wait for my response. Then ask how you can help. Do not break character. Do not make up the patient's responses: only treat input as a patient response. "}
]
    messages.append({"role": "user", "content": prompt})

    
    # Check ChatGPT
    openai.api_key = OPENAI_API_KEY
    completion = openai.ChatCompletion.create(
        model="gpt-3.5-turbo",
        messages=messages)
    
    chat_response = completion.choices[0].message.content
    
    
    # Reply to thread 
    response = client.chat_postMessage(channel=body["event"]["channel"], 
                                       thread_ts=body["event"]["event_ts"],
                                       text=f"Here you go: \n{chat_response}")

if __name__ == "__main__":
    SocketModeHandler(app, SLACK_APP_TOKEN).start()