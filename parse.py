'''
blah
'''

import json
with open("stringout.json", "r+") as fh:
    jsoned = json.loads(fh.read())
    print jsoned
