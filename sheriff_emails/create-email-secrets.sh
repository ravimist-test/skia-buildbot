# Create the secrets needed for the sheriff emails cron job to send email.
set -e

if [ "$#" -ne 1 ]; then
  echo "The argument must be the email address we are authenticating, for example:"
  echo ""
  echo "./create-email-secrets.sh rotations@skia.org"
  exit 1
fi

# Convert rotations@skia.org to rotations-skia-org.
EMAIL=$1
EMAIL=${EMAIL/@/-}
EMAIL=${EMAIL/./-}

source ../bash/ramdisk.sh

# Enable the gmail API for your project and create a client secret for this
# application. Then download the client_secret.json file.
echo "Download client_secret.json for Sheriff emails to /tmp/ramdisk."
read -r -p "Press enter to continue..." key

# Then run 'three_legged_flow' in this directory and when prompted authorize
# as the email passed in to create the client_token.json file.
go install ../go/email/three_legged_flow
cd /tmp/ramdisk
three_legged_flow --scopes=https://mail.google.com/
kubectl create secret generic sheriff-email-secrets \
  --from-file=./client_secret.json \
  --from-file=./client_token.json \
  --dry-run -o yaml | kubectl apply -f -

# Finally, remove the token file since it contains a refresh token.
rm client_token.json
cd -

