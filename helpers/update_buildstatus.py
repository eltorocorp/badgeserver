
import argparse
import urllib.request
import urllib.parse


def update_buildstatus():
    """
    update_buildstatus sends a POST request to the badge server to create/update a build status badge.

    Arguments:
    - projectname
    - status [building|failing|passing]

    If the build status is building, the resulting badge will be blue.
    If the build status is failing, the resulting badge will be red.
    If the build status is passing, the resulting badge will be green.
    If a value other than building, failing, or passing is supplied, the script will use that value with a blue badge.

    Example: python update_buildstatus --projectname foo --status passing
    """
    args = get_args()
    template = 'http://badgeserver.mydomain.com?project={}&item={}&value={}&color={}'
    color = "blue"
    if args.status == "failing":
        color = "red"
    if args.status == "passing":
        color = "green"
    url = template.format(args.projectname, "build", args.status, color)
    data = urllib.parse.urlencode({}).encode("utf-8")
    urllib.request.urlopen(url, data=data)


def get_args():
    parser = argparse.ArgumentParser()
    parser.add_argument('--projectname', required=True)
    parser.add_argument('--status', required=True)
    return parser.parse_args()


if __name__ == "__main__":
    update_buildstatus()
