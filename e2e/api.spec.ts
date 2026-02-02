import { test, expect } from '@playwright/test';

test.describe('Scrumer Backend E2E', () => {
  test('GET /ping returns pong', async ({ request }) => {
    const res = await request.get('/ping');
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body).toEqual({ message: 'pong' });
  });

  test('GET /graphql serves GraphiQL', async ({ request }) => {
    const res = await request.get('/graphql', {
      headers: { Accept: 'text/html' },
    });
    expect(res.ok()).toBeTruthy();
    const text = await res.text();
    expect(text).toContain('GraphiQL');
  });

  test('POST /graphql hello query returns world', async ({ request }) => {
    const res = await request.post('/graphql', {
      data: { query: 'query { hello }' },
      headers: { 'Content-Type': 'application/json' },
    });
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body.data?.hello).toBe('world');
  });

  test('POST /graphql createUser mutation', async ({ request }) => {
    const username = `e2euser_${Date.now()}`;
    const email = `e2e_${Date.now()}@example.com`;
    const res = await request.post('/graphql', {
      data: {
        query: `mutation CreateUser($u: String!, $e: String!, $p: String!) {
          createUser(username: $u, email: $e, password: $p) { id username email }
        }`,
        variables: {
          u: username,
          e: email,
          p: 'password123',
        },
      },
      headers: { 'Content-Type': 'application/json' },
    });
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body.errors).toBeFalsy();
    expect(body.data?.createUser?.username).toBe(username);
  });
});
