/// This is a global variable that keeps track of how many request we have going on in htmx.
// for whatever reason I was getting more than one trigger of the "onLoad" for htmx
// when trying to do some after load manipulations. Specifically, I was getting more than one
// notification shown because for some reason it was triggered twice. This is a hleper
// to keep track of that and not show more than just one notification per requets.
//
// This is probably stupid and there is a better way to do this, but for now, we just relallllly dont' care.
let activeRequests = 0;

document.addEventListener("DOMContentLoaded", () => {
  document.body.addEventListener("htmx:configRequest", () => {
    activeRequests++;
  });

  document.body.addEventListener("htmx:afterRequest", handleHtmxAfterRequest);
});

/// Function to handle the notifications coming form the htmx response that are
//stored in the headres. We want to read it and then show it to the user with the UI kit stuff.
function handleNotifications(event) {
  const xhr = event.detail.xhr;

  const errorsHeader = xhr.getResponseHeader("X-Errors");
  const notificationsHeader = xhr.getResponseHeader("X-Notifications");

  if (errorsHeader) {
    try {
      const errors = JSON.parse(errorsHeader);
      console.log("ERRORS:: ", errors);
      errors.forEach((m) => {
        UIkit.notification(m, {
          pos: "bottom-center",
          status: "destructive",
          timeout: 2500,
        });
      });
    } catch (e) {
      console.error(e);
    }
  }

  if (notificationsHeader) {
    try {
      const notifications = JSON.parse(notificationsHeader);
      console.log("NOTIFICATIONS:: ", notifications);
      notifications.forEach((m) => {
        UIkit.notification(m, {
          pos: "bottom-center",
          status: "success",
          timeout: 2500,
        });
      });
    } catch (e) {
      console.error(e);
    }
  }
}

// Helper to be able to "remove" the modal from the DOM when the modal is closed.
// This is specifically done for the htmx response with the modal.
//There are multiple ways how the modal can be closed and this is hard set up with the
//hyperscirpt. Or at least skill issue and I don't kno how.
function handleModalRemoval(event) {
  const response = event.detail.xhr.response;
  const isModal = response.search("uk-modal");
  if (isModal > -1) {
    const firstId = response.match(/id="([^"]+)"/)[1];
    const modal = document.querySelector(`#${firstId}`);
    UIkit.util.on(modal, "hidden", () => {
      modal.remove();
    });
  }
}
// This is a hnadler of the htmx request. I atttach all kinds of listeners and extra
// behavior for the newly created contet from htmx.
function handleHtmxAfterRequest(event) {
  if (activeRequests === 0) {
    return;
  }
  activeRequests--;
  handleNotifications(event);
  handleModalRemoval(event);
}
