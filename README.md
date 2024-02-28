# Motmot

## Running

### Backend

Currently, the database and backend API are configured to run with Docker. This is done to simplify configuration and networking.

Use this command to start the machines:

```bash
# We want to spin UP the machines
# The -d flag "detaches" you from the machine, otherwise you'd need to keep the terminal open
docker compose up -d
```

In order to stop the machines (otherwise they will continue running even after closing your terminal), use this command:

```bash
# We want to spin DOWN the machines
docker compose down
```
### Frontend

The frontend resides in the "www" folder. In order to run the front end, make sure you have Node.js installed and run the following command.

```bash
npm install # Only needed the first time
npm run dev
```

Then go to the url printed in the terminal.

### Usage

Access the API documentation at `http://localhost:3000/api/v1/docs/index.html`.

## Resources

- [Tailwind CSS](https://tailwindcss.com/docs/flex-basis)
- [daisyUI](https://daisyui.com/components/)
- [SvelteKit Docs](https://kit.svelte.dev/docs/routing)