function initKeycloak() {
    var keycloak = new Keycloak();
    keycloak.init().then(function (authenticated) {
        alert(authenticated ? 'authenticated' : 'not authenticated');
    }).catch(function () {
        alert('failed to initialize');
    });
}

initKeycloak();