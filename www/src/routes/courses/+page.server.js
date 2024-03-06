// This fetches the URL of the backend API from environment variables.
// The actual value is specified in the file atmotmot/www/.env
// It's good practice not to hardcode this since it can change in the
// real world.
import { BASE_API_PATH } from '$lib/env';

/** @type {import('./$types').PageServerLoad} */
export async function load({ fetch, locals }) {
    // This function is called on the server before the page is sent to the client
    // If the user has not logged in, we just want to list all the courses at UF;
    // otherwise, we want to know which groups the user is in by querying the API.

    // If this variable is never populated it'll be undefined, which will tell
    // our webpage that no user has logged in.
    let data;

    // If the "token" local has been set (which stores the current user's token),
    // then we know a user has logged in, so we need to check which groups they
    // are in.
    if(locals.token !== undefined)
    {
        // To see which endpoint to call, we can look at the debug page for the
        // API (see the README).

        // Notice the 'async' tag on the function declaration and the 'await'
        // keyword on this call. Querying an API takes time, noticeable time,
        // so we can't just freeze the entire UI while we're doing it.
        // async and await tell JavaScript that until this function is resolved,
        // go back to whatever you were doing before this function was called.
        const res = await fetch(`${BASE_API_PATH}/user/groups`, {
            // The API's debug page labels this endpoint as a GET request
            method: 'get',

            // In order to authenticate with the backend and ask for the groups,
            // we need to send the token along with our request as proof that we
            // are who we say we are.
            // In our case, the token is included automatically when we add this line
            credentials: 'include',

            // In order to allow us to include credentials, we need to enable CORS.
            // This is a bunch of bullshit that I'm not going to explain in comments
            // so ask me if you're curious.
            mode: 'cors'
        });

        // All of the backend API endpoints return JSON (look it up if you aren't
        // familiar with it). We also need to add await here too, since at this
        // point we haven't actually received the request's data.
        data = await res.json();
    }

    // Return this data to the webpage under the label 'groups'
    return {
        groups: data
    }
}