import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart

# Set up email
msg = MIMEMultipart()
msg['From'] = 'test@example.com'
msg['To'] = 'bold-quick41@localhost.local'
msg['Subject'] = 'Python Test Email'
msg.attach(MIMEText('This is a test email sent from Python', 'plain'))

# Send the email
server = smtplib.SMTP('localhost', 2525)
server.send_message(msg)
server.quit()

print("Email sent!")