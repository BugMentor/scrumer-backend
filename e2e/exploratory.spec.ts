import { test, expect } from '@playwright/test';

test.describe('Exploratory UI', () => {
  test('GraphiQL page is served and contains UI', async ({ request }) => {
    const res = await request.get('/graphql', {
      headers: { Accept: 'text/html' },
    });
    expect(res.ok()).toBeTruthy();
    const html = await res.text();
    expect(html).toContain('GraphiQL');
    expect(html).toMatch(/graphql|query/i);
  });
});
