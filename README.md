# Coach Personnel IA

Application web de génération de programmes d'entraînement personnalisés par intelligence artificielle (Claude d'Anthropic).

## Fonctionnalités

- **Génération de programme** — Créez un programme d'entraînement sur mesure en indiquant votre profil (niveau, objectifs, équipement, fréquence)
- **Affichage structuré** — Chaque séance présente le mouvement, les séries/répétitions, l'intensité (%RM ou RPE), le tempo et les consignes techniques
- **Adaptation au ressenti** — Renseignez votre ressenti de la semaine précédente pour que le programme s'ajuste en conséquence
- **Téléchargement PDF** — Exportez n'importe quelle séance en PDF formaté
- **Timer d'entraînement** — Minuteur interactif adapté à la séance du jour (phases de travail et de repos, bip sonore)
- **Historique** — Tous les programmes générés sont sauvegardés en base SQLite
- **Authentification** — Accès protégé par mot de passe configurable
- **Serveur MCP** — Intégration avec Claude Desktop / Claude Code pour générer des programmes depuis le chat

---

## Démarrage rapide avec Docker

```bash
# 1. Cloner le dépôt
git clone <repo-url>
cd personal-coach

# 2. Configurer l'environnement
cp .env.example .env
# Éditez .env et renseignez votre clé API Anthropic
# ANTHROPIC_API_KEY=sk-ant-...

# 3. Lancer l'application
docker-compose up -d

# 4. Ouvrir dans le navigateur
# http://localhost:8080
# Mot de passe par défaut : coach2024
```

---

## Installation manuelle

### Prérequis

| Outil | Version minimale |
|-------|-----------------|
| Go    | 1.22+           |
| Node  | 18+             |
| npm   | 9+              |

### Construction

```bash
# 1. Configurer l'environnement
cp .env.example .env
# Éditez .env : ANTHROPIC_API_KEY, APP_PASSWORD, PORT...

# 2. Construire (frontend → backend)
./scripts/build.sh

# 3. Lancer le serveur
source .env && ./personal-coach
# Ou directement :
ANTHROPIC_API_KEY=sk-ant-... APP_PASSWORD=monmotdepasse ./personal-coach
```

Le serveur démarre sur le port `8080` (modifiable via `PORT`).
Ouvrez **http://localhost:8080** et connectez-vous avec votre mot de passe.

---

## Variables d'environnement

| Variable           | Défaut        | Description                                      |
|--------------------|---------------|--------------------------------------------------|
| `ANTHROPIC_API_KEY`| *obligatoire* | Clé API Anthropic (Claude)                       |
| `APP_PASSWORD`     | `coach2024`   | Mot de passe de l'application                    |
| `PORT`             | `8080`        | Port HTTP d'écoute                               |
| `DATA_DIR`         | `./data`      | Répertoire de la base de données SQLite          |

---

## Structure du projet

```
personal-coach/
├── backend/                  # Serveur Go (Gin)
│   ├── main.go               # Point d'entrée (HTTP ou MCP)
│   ├── database/             # SQLite : migrations et CRUD
│   │   ├── db.go
│   │   ├── migrations.go     # Migrations versionnées (ne jamais modifier)
│   │   └── store.go
│   ├── handlers/             # Contrôleurs REST
│   │   ├── auth.go           # Login / logout / session
│   │   └── program.go        # Génération, PDF, timer, liste
│   ├── mcp/
│   │   └── server.go         # Serveur MCP (JSON-RPC 2.0 stdio)
│   ├── models/
│   │   └── models.go         # Structures de données
│   └── services/
│       ├── claude.go         # Intégration SDK Claude
│       ├── pdf.go            # Génération PDF (gofpdf)
│       └── timer.go          # Construction de la séquence timer
├── frontend/                 # Vue 3 + Tailwind CSS (Vite)
│   └── src/
│       ├── views/
│       │   ├── LoginView.vue
│       │   ├── HomeView.vue        # Formulaire de création
│       │   ├── ProgramView.vue     # Affichage du programme
│       │   └── ProgramsListView.vue # Historique
│       ├── components/
│       │   └── TimerModal.vue      # Timer interactif
│       └── stores/
│           ├── auth.js             # Store Pinia : authentification
│           └── program.js          # Store Pinia : programmes
├── scripts/
│   └── build.sh              # Script de compilation complet
├── Dockerfile                # Multi-stage : Node → Go → Alpine
├── docker-compose.yml
└── .env.example
```

---

## API REST

Toutes les routes `/api/*` nécessitent une session valide (cookie `coach_session`).

| Méthode | Route                             | Description                        |
|---------|-----------------------------------|------------------------------------|
| `POST`  | `/auth/login`                     | Connexion (body: `{"password":"…"}`)|
| `POST`  | `/auth/logout`                    | Déconnexion                        |
| `GET`   | `/auth/status`                    | Vérification de session            |
| `GET`   | `/health`                         | Health check                       |
| `POST`  | `/api/programs/generate`          | Générer un programme               |
| `GET`   | `/api/programs`                   | Lister tous les programmes         |
| `GET`   | `/api/programs/:id`               | Récupérer un programme             |
| `GET`   | `/api/programs/:id/pdf`           | Télécharger en PDF                 |
| `GET`   | `/api/programs/:id/timer/:day`    | Timer de la séance (0-based)       |

### Exemple : générer un programme

```bash
# 1. Se connecter
curl -c cookies.txt -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"password":"coach2024"}'

# 2. Générer un programme
curl -b cookies.txt -X POST http://localhost:8080/api/programs/generate \
  -H "Content-Type: application/json" \
  -d '{
    "person": {
      "name": "Alice",
      "age": 30,
      "weight": 65,
      "height": 168,
      "level": "intermediate",
      "goals": ["muscle_gain", "strength"],
      "equipment": ["barbell", "dumbbell"]
    },
    "days_per_week": 4,
    "weeks": 6
  }'
```

### Corps de la requête `/api/programs/generate`

```json
{
  "person": {
    "name": "Alice",
    "age": 30,
    "weight": 65.0,
    "height": 168.0,
    "level": "beginner | intermediate | advanced",
    "goals": ["weight_loss", "muscle_gain", "strength", "endurance", "flexibility", "general_fitness"],
    "equipment": ["bodyweight", "dumbbell", "barbell", "machine", "kettlebell", "bands", "pullup_bar"]
  },
  "days_per_week": 3,
  "weeks": 4,
  "feedback": {
    "energy_level": 7,
    "soreness_level": 3,
    "motivation_level": 8,
    "completed_days": 3,
    "notes": "Légère douleur à l'épaule droite"
  }
}
```

---

## Serveur MCP

Intégrez l'application dans **Claude Desktop** ou **Claude Code** pour générer des programmes directement depuis le chat.

### Configuration Claude Desktop

Ajoutez dans `claude_desktop_config.json` :

```json
{
  "mcpServers": {
    "personal-coach": {
      "command": "/chemin/vers/personal-coach",
      "args": ["mcp"],
      "env": {
        "ANTHROPIC_API_KEY": "sk-ant-..."
      }
    }
  }
}
```

### Outils MCP disponibles

| Outil                      | Description                                          |
|----------------------------|------------------------------------------------------|
| `generate_workout_program` | Génère un programme complet selon le profil          |
| `get_workout_timer`        | Construit la séquence timer à partir d'un programme  |

---

## Développement

### Backend seul (avec hot-reload via air ou go run)

```bash
cd backend

# Copier le frontend compilé (nécessaire pour go:embed)
npm --prefix ../frontend run build && cp -r ../frontend/dist ./dist

# Lancer
GOROOT=/home/banux/go GOPATH=/home/banux/go \
  ANTHROPIC_API_KEY=sk-ant-... \
  go run main.go
```

### Frontend en mode développement

```bash
cd frontend
npm run dev
# Disponible sur http://localhost:5173
# Proxie les requêtes /api/* vers http://localhost:8080
```

Pour que le proxy fonctionne en dev, ajoutez dans `vite.config.js` :

```js
server: {
  proxy: {
    '/api': 'http://localhost:8080',
    '/auth': 'http://localhost:8080',
  }
}
```

### Tests

```bash
cd backend
GOROOT=/home/banux/go GOPATH=/home/banux/go go test ./database/... -v
# 6 tests : CRUD programmes + idempotence des migrations
```

---

## Docker

```bash
# Démarrer
docker-compose up -d

# Voir les logs
docker-compose logs -f

# Reconstruire après modifications
docker-compose up -d --build

# Arrêter
docker-compose down

# Supprimer la base de données (volume)
docker-compose down -v
```

Les données SQLite sont persistées dans le volume Docker `coach-data`.

---

## Licence

Projet personnel — usage libre.
