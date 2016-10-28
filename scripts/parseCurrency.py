import urllib.request
import xml.etree.cElementTree

results = {}
resp = urllib.request.urlopen("http://www.currency-iso.org/dam/downloads/lists/list_one.xml").read().decode('utf-8')
e = xml.etree.cElementTree.fromstring(resp)
for entry in e[0].getchildren():
    code = entry.find('Ccy')
    if code is not None:
        if entry.find('CcyMnrUnts').text != 'N.A.':
            results[code.text] = {
                'name': entry.find('CcyNm').text,
                'digits': entry.find('CcyMnrUnts').text,
            }

for code, d in results.items():
    print(
'''"%s": currency{
    Name: "%s",
    Code: "%s",
    DigitsAfterDecimal: %s,
},''' % (code, d['name'], code, d['digits']))
