import requests

def handle(req):
    """handle a request to the function
    Args:
        req (str): request body
    """

    return requests.get("https://wttr.in/" + req).text

