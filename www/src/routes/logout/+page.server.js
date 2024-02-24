import { redirect } from '@sveltejs/kit';

/** @type {import('./$types').Actions} */
export const actions = {
    default: async ({ cookies, url }) => {
        cookies.delete("token", { path: '/' });
        redirect(303, url.searchParams.get('redirectTo') ?? '/');
    }
}