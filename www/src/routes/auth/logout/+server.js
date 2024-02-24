import { BASE_API_PATH } from '$lib/env';
import { parse } from 'cookie';

/** @type {import('./$types').RequestHandler} */
export async function POST({ cookies }) {
    cookies.delete('token', { path: '/' });
    return new Response();
}