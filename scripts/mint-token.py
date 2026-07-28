import base64, hashlib, hmac, json, os, sys, time

def b64(b): return base64.urlsafe_b64encode(b).rstrip(b'=')

secret = os.environ['JWT_SECRET'].encode()
sub, phone, role = sys.argv[1], sys.argv[2], sys.argv[3]
now = int(time.time())
hdr = b64(json.dumps({"alg":"HS256","typ":"JWT"},separators=(',',':')).encode())
pl  = b64(json.dumps({"sub":sub,"phone":phone,"role":role,"is_super_admin":False,
                      "iss":"urja","iat":now,"exp":now+14400},separators=(',',':')).encode())
msg = hdr + b'.' + pl
sig = b64(hmac.new(secret, msg, hashlib.sha256).digest())
sys.stdout.write((msg + b'.' + sig).decode())
