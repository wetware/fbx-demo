from tikapi import TikAPI, ValidationException, ResponseException
import os

from tikapi.api import APIResponse


# Notification format: https://tikapi.io/documentation/#tag/Profile/operation/user.notifications
# Additional context for the comments: https://tikapi.io/documentation/#tag/Public/operation/public.commentRepliesList


api_key: str = os.getenv("API_KEY", "")
account_key: str = os.getenv("ACCOUNT_KEY", "")


if api_key == "" or account_key == "":
    raise ValueError("api_key and/or account_key is not set")


api = TikAPI(api_key)
User = api.user(accountKey=account_key)


def notification_loop():
    while True:
        try:
            response = User.notifications(filter="mentions")
            while response:
                response = process_notification(response)
        except ValidationException as e:
            print(e, e.field)

        except ResponseException as e:
            print(e, e.response.status_code)


def process_notification(response: APIResponse) -> APIResponse:
    # TODO: extract comment body and send it to Wetware agent.
    min_time = response.json().get("min_time")  # TODO: change field
    return response.next_items()


if __name__ == "__main__":
    notification_loop()
