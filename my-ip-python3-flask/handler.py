import requests

def handle(req):
    """handle a request to the function
    Args:
        req (str): request body
    """

    return requests.get("http://ip.jsontest.com/").text
