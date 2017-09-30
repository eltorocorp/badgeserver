import argparse
import urllib.request
import urllib.parse


def update_version():
    """
    update_version sends a POST request to the badge server to create/update a version badge.

    Arguments:
    - projectname
    - version

    The version badge will be blue.

    Example: python update_buildstatus --projectname foo --version 1.0.0.0
    """
    args = get_args()
    template = 'http://badgeserver.mydomain.com?project={}&item=version&value={}&color=blue'
    url = template.format(args.projectname, args.version)
    data = urllib.parse.urlencode({}).encode("utf-8")
    urllib.request.urlopen(url, data=data)


def get_args():
    parser = argparse.ArgumentParser()
    parser.add_argument('--projectname', required=True)
    parser.add_argument('--version', required=True)
    return parser.parse_args()


if __name__ == "__main__":
    update_version()
