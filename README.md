## Api qui permet de controller des télécommandes de volets

### Installation

<p align="center">
    <img src="example.jpg" width="512" alt="Installation"/>
</p>

Les télécommandes ayant un fil commun entre les deux boutons je connecte uniquement 3 fils.

### Intégration dans home assistant

Tout d'abord il faut créer une "commande" qui permettera de contacter l'api

Dans le fichier configuration.yaml ajouter ceci :

```
rest_command:                    
  shutter:                                                       
    url: "http://ip:port/?s={{name}}&p={{position}}"
    headers:                  
      Authorization: !secret shutters_secret
```

La variable name sert a déterminer le volet et position permet d'envoyer soit `up` ou `down`

le token d'accès lui est stocké dans le fichier secrets.yaml :

``shutters_secret: "Bearer example"``

Ensuite il faut ajouter un nouveau volet toujours dans configuration.yaml :

```
cover:
  - platform: template
    covers:
      shutter_example:
        device_class: shutter
        friendly_name: "Volet Exemple"
        open_cover:
          service: rest_command.shutter
          data:
            name: "example"
            position: "up"
        close_cover:
          service: rest_command.shutter
          data:
            name: "example"
            position: "down"
```