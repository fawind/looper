import sys
import logging
from flask import Flask

app = Flask(__name__)

logging.basicConfig(
    format="[%(asctime)s|%(name)s-%(funcName)s(%(lineno)d)|%(levelname)s]: %(message)s",
    level="INFO",
    stream=sys.stdout,
)

log = logging.getLogger(__file__)


@app.route("/sink/<int:request_id>")
def registered_workers(request_id: int) -> str:
    log.info("Received %s", request_id)
    return "200"


if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0")
