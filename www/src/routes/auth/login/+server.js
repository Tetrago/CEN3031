import { BASE_API_PATH } from '$lib/env';
import { parse } from 'cookie';

/** @type {import('./$types').RequestHandler} */
export async function POST({ cookies, fetch, request }) {
    const res = await fetch(`${BASE_API_PATH}/auth/login`, {
        method: 'post',
        mode: 'cors',
        body: JSON.stringify(await request.json())
    });

    if(res.ok) {
        cookies.set('token', parse(res.headers.getSetCookie()[0]).token, { path: '/' });
    }

    return new Response(JSON.stringify(await res.json()), {
        status: res.status
    });
}