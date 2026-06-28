db = db.getSiblingDB("inventory");

const parts = [
  {
    _id: "11111111-1111-1111-1111-111111111111",
    name: "Rocket Engine",
    description: "Main rocket engine",
    price: 150000,
    stock_quantity: 10,
    category: "Engine",
    manufacturer_name: "Space Factory",
    manufacturer_country: "Russia",
    manufacturer_website: "https://space-factory.example.com",
    length: 4.5,
    width: 2.5,
    height: 3,
    weight: 1200,
    tags: ["engine", "rocket"],
    created_at: new Date("2026-01-01T00:00:00Z"),
    updated_at: new Date("2026-01-01T00:00:00Z")
  },
  {
    _id: "22222222-2222-2222-2222-222222222222",
    name: "Fuel Tank",
    description: "Reinforced spacecraft fuel tank",
    price: 42000,
    stock_quantity: 20,
    category: "Fuel",
    manufacturer_name: "Orbital Systems",
    manufacturer_country: "USA",
    manufacturer_website: "https://orbital-systems.example.com",
    length: 6,
    width: 2.8,
    height: 2.8,
    weight: 850,
    tags: ["fuel", "tank"],
    created_at: new Date("2026-01-01T00:00:00Z"),
    updated_at: new Date("2026-01-01T00:00:00Z")
  },
  {
    _id: "33333333-3333-3333-3333-333333333333",
    name: "Space Porthole",
    description: "Radiation-resistant spacecraft porthole",
    price: 5000,
    stock_quantity: 25,
    category: "Porthole",
    manufacturer_name: "Space Glass",
    manufacturer_country: "Germany",
    manufacturer_website: "https://space-glass.example.com",
    length: 0.1,
    width: 1,
    height: 1,
    weight: 45,
    tags: ["glass", "porthole"],
    created_at: new Date("2026-01-01T00:00:00Z"),
    updated_at: new Date("2026-01-01T00:00:00Z")
  },
  {
    _id: "44444444-4444-4444-4444-444444444444",
    name: "Stabilizer Wing",
    description: "Composite stabilizer wing",
    price: 28000,
    stock_quantity: 16,
    category: "Wing",
    manufacturer_name: "Aero Space",
    manufacturer_country: "France",
    manufacturer_website: "https://aero-space.example.com",
    length: 7.2,
    width: 3.5,
    height: 0.4,
    weight: 310,
    tags: ["wing", "stabilizer"],
    created_at: new Date("2026-01-01T00:00:00Z"),
    updated_at: new Date("2026-01-01T00:00:00Z")
  }
];

for (const part of parts) {
  const fields = { ...part };
  delete fields._id;

  db.parts.updateOne(
    { _id: part._id },
    { $setOnInsert: fields },
    { upsert: true }
  );
}
