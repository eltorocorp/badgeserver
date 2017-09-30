import argparse
import urllib.request
import urllib.parse


def update_coverage():
    """
    update_coverage sends a POST request to the badge server to create/update a test coverage status badge.

    Arguments:
    - projectname
    - coverage
    - standard

    If the supplied coverage value is less than the standard, a red badge is used, otherwise a green badge is used.

    Example: python update_buildstatus --projectname foo --coverage 55.2 --standard 100.0
    """
    args = get_args()
    template = 'http://badgeserver.mydomain.com?project={}&item={}&value={}&color={}'
    color = "green" if args.coverage >= args.standard else "red"
    value = str(args.coverage) + '%25'
    url = template.format(args.projectname, "coverage", value, color)
    data = urllib.parse.urlencode({}).encode("utf-8")
    urllib.request.urlopen(url, data=data)


def get_args():
    parser = argparse.ArgumentParser()
    parser.add_argument('--projectname', required=True)
    parser.add_argument('--coverage', required=True, type=float)
    parser.add_argument('--standard', required=True, type=float)
    return parser.parse_args()


if __name__ == "__main__":
    update_coverage()
