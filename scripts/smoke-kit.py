"""Run against the disposable CI image with a fresh database."""
import os
from playwright.sync_api import sync_playwright, expect

with sync_playwright() as pw:
    browser = pw.chromium.launch(channel="chrome")
    page = browser.new_page(base_url="http://127.0.0.1:9000")
    # Direct entry must serve the SPA, without redirecting to the site root.
    page.goto("/login")
    page.locator('[autocomplete="username"]').fill("admin")
    page.locator('[autocomplete="current-password"]').fill(os.environ["MEMO_ADMIN_PASSWORD"])
    page.get_by_role("button", name="登录").click()
    composer = page.get_by_placeholder("现在的想法是...")
    composer.fill("UpgradeProbe")
    with page.expect_response(lambda r: r.url.endswith("/api/memos") and r.request.method == "POST") as created:
        page.locator('button[title^="发送"]').click()
    assert created.value.ok
    search = page.locator('input[placeholder*="搜索"]').first
    with page.expect_response(lambda r: "/api/memos?" in r.url and "UpgradeProbe" in r.url) as result:
        search.fill("UpgradeProbe")
    assert result.value.ok and "UpgradeProbe" in result.value.text()
    page.goto("/notebooks")
    page.reload()
    expect(page).to_have_url("http://127.0.0.1:9000/notebooks")
    page.goto("/resources")
    with page.expect_response(lambda r: r.url.endswith("/api/resources") and r.request.method == "POST") as uploaded:
        page.locator('input[type=file]').set_input_files({"name": "probe.txt", "mimeType": "text/plain", "buffer": b"smoke"})
    assert uploaded.value.status == 201
    page.evaluate("localStorage.setItem('token', 'invalid')")
    page.reload()
    page.wait_for_url("**/login")
    assert page.evaluate("localStorage.getItem('token')") is None
    browser.close()
