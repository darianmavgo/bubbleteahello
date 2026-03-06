#!/bin/bash

# Step 1: Auth Check (like Bubble Tea's Init())
if ! npx wrangler whoami >/dev/null 2>&1; then
  echo "🔐 Logging in to Cloudflare—grab your token!"
  npx wrangler login
fi
echo "✅ Auth good—let's redirect like a pro."

# Create Project (Step 2 Starter)
npx wrangler generate grape-redirect-worker --type javascript --yes
cd grape-redirect-worker

# Write the Redirect Worker Code (Step 2 Core + Step 3 Edges)
cat > src/index.js << 'EOF'
export default {
  async fetch(request) {
    const url = new URL(request.url);
    let targetHost = '2122grapeave.darianhickman.com';

    // Handle WWW & HTTP edges: Normalize to HTTPS non-www
    if (url.hostname === 'www.2122grapeave.mavgo.com') {
      url.hostname = '2122grapeave.mavgo.com';  // First to base for logic
    }
    if (url.protocol === 'http:') {
      url.protocol = 'https:';
    }

    // Core Redirect: Mavgo → Darian (preserve path/query)
    if (url.hostname === '2122grapeave.mavgo.com' || url.hostname === 'www.2122grapeave.mavgo.com') {
      const targetUrl = `https://${targetHost}${url.pathname}${url.search}`;
      return Response.redirect(targetUrl, 301);  // Permanent, SEO-friendly
    }

    // Fallback: Pass through if not matching (safety net)
    return fetch(request);
  },
};
EOF

# Config: Bind Route to Subdomain (Step 2 Polish)
cat > wrangler.toml << 'EOF'
name = "grape-redirect-worker"
main = "src/index.js"
compatibility_date = "2025-11-27"

[[routes]]
pattern = "2122grapeave.mavgo.com/*"
zone_name = "mavgo.com"

[[routes]]
pattern = "www.2122grapeave.mavgo.com/*"
zone_name = "mavgo.com"
EOF

# Deploy It (Steps 2-3 Done!)
echo "🚀 Deploying redirect Worker..."
npx wrangler deploy
WORKER_URL=$(npx wrangler --version | grep -o 'workers.dev' | head -1)  # Nah, just confirm deploy output

# Step 4: Test the Glow-Up
echo "🧪 Testing redirect..."
sleep 10  # Prop wait
if curl -I -s https://2122grapeave.mavgo.com/ | grep -q "301"; then
  echo "✅ Redirect LIVE! Old → New Grape Ave page."
  curl -s https://2122grapeave.mavgo.com/ | head -1  # Peek final land
else
  echo "❌ Test flopped—check deploy logs: npx wrangler tail"
fi

# Cleanup? Uncomment if temp: cd .. && rm -rf grape-redirect-worker
echo "🎉 Script done—your old site's a forwarding ghost. Prod URL: $(npx wrangler deploy --dry-run | grep 'workers.dev')"
