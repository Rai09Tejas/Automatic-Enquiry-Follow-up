# Golf Club Enquiry Follow-Up

This is a small enquiry form for membership, society days, and events. Visitors submit an enquiry, and staff can view enquiries, mark them as contacted, and see which new enquiries need a follow-up.

## Running it locally

You need Go 1.22+ and Node.js.

Build the React client:

```bash
cd client
npm install
npm run build
```

Run `npm run dev` inside `client/` and use the URL it displays.

------------------

Run the Go server:

```bash
cd ../server
go run .
```

Open [http://localhost:8080](http://localhost:8080). The staff page is at `/staff`. The SQLite database is created in `server/enquiries.db`.


## Main decisions

- React and Vite to keep the client simple and easy to edit.
- Go provides the API and serves the built client.
- SQLite keeps setup small, now there is no separate database server.
- New enquiries older than 48 hours are marked as due for follow-up when they are listed, this is calculated real-time rather than stored.
- Confirmation emails are currently logged in the console in the server terminal instead of being sent.

## Adding it to an existing club website

- The club’s web person would normally deploy this app separately, for example on a small hosting service or the club’s existing server. They would build the React client, run the Go server, and give it a public address such as `enquiries.clubname.co.uk`.

They could then add a clear “Make an enquiry” button to the club website that links to that address. If the hosting setup supports it, they could also place the form inside a page using an iframe. A link is usually simpler and more reliable for WordPress or Wix, while an iframe keeps visitors visually closer to the main website.

- The form should use the club’s logo, colours, font choices, and wording. Its layout is already deliberately simple, and the web person could adjust the React stylesheet to match the club’s design. The form should be checked on a phone before launch, if it is dynamic or not.

- For enquirers, only the necessary fields should be requested: name, email, enquiry type, and message. To reduce spam, the live version should add a CAPTCHA, server-side validation, rate limiting, and possibly email verification. Staff pages should also require a proper login before the app is used publicly.

- If the club redesigns its website later, the form does not need to be rebuilt from scratch. The web person can update the client colours, logo, spacing, and wording, then rebuild and redeploy it. If the form remains at the same public address, the button or link on the redesigned website should continue to work. An iframe may need its width and height checked again after a redesign, which is another reason a normal link is often preferable.

## What I would build next

- Staff login and permissions (with role-based access control).
- Real confirmation emails (maybe with Go's builtin lib net/smtp) with send and failure tracking.
- Spam protection and rate limiting.
- Search and filters for staff (by name, type of enquiry, created date, status and followup due status).
- Automated tests for validation and the 48-hour follow-up rule.
- A production database and backups instead of a local SQLite file.
- Better UI😄

## AI tools used

OpenAI Codex was used as a coding assistant. It helped reorganise the project into separate client and server folders, simplify the Go code, update the server setup, run build checks, and paraphrase this documentation. 
