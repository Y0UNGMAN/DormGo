from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class GoAPIConfig:
    api_url: str
    api_key: str
    model: str


def load_go_api_config(path: str = None) -> GoAPIConfig:
    config_path = Path(path) if path else _default_config_path()
    values = _read_api_section(config_path)
    missing = [key for key in ("ApiUrl", "ApiKey", "ApiModel") if not values.get(key)]
    if missing:
        raise RuntimeError(f"Go config api section missing: {', '.join(missing)}")
    return GoAPIConfig(
        api_url=values["ApiUrl"],
        api_key=values["ApiKey"],
        model=values["ApiModel"],
    )


def _default_config_path() -> Path:
    return Path(__file__).resolve().parents[1] / "backend" / "config.yaml"


def _read_api_section(config_path: Path) -> dict:
    if not config_path.exists():
        raise RuntimeError(f"Go config not found: {config_path}")

    values = {}
    in_api = False
    for raw_line in config_path.read_text(encoding="utf-8").splitlines():
        line = raw_line.rstrip()
        stripped = line.strip()
        if not stripped or stripped.startswith("#"):
            continue
        if not line.startswith(" ") and stripped.endswith(":"):
            in_api = stripped[:-1] == "api"
            continue
        if in_api and ":" in stripped:
            key, value = stripped.split(":", 1)
            values[key.strip()] = value.strip().strip('"').strip("'")
    return values
