import { StatusBar } from "expo-status-bar";
import { useEffect, useState } from "react";
import {
  ActivityIndicator,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
  Platform,
} from "react-native";

const API_URL =
  Platform.OS === "web"
    ? process.env.EXPO_PUBLIC_API_URL ?? "http://localhost:8080"
    : process.env.EXPO_PUBLIC_API_URL ?? "http://10.0.2.2:8080";

export default function App() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [token, setToken] = useState<string | null>(null);
  const [name, setName] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  // When we have a token, query the me endpoint
  useEffect(() => {
    if (!token) return;

    setLoading(true);
    setError(null);

    fetch(`${API_URL}/graphql`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        query: `query { me { name } }`,
      }),
    })
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => {
        if (data.errors) throw new Error(data.errors[0].message);
        setName(data.data.me.name);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, [token]);

  const handleLogin = async () => {
    setLoading(true);
    setError(null);

    try {
      const res = await fetch(`${API_URL}/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || `HTTP ${res.status}`);
      }

      setToken(data.token);
    } catch (err: any) {
      setError(err.message);
      setLoading(false);
    }
  };

  const handleLogout = () => {
    setToken(null);
    setName(null);
    setEmail("");
    setPassword("");
  };

  // Logged in view
  if (name) {
    return (
      <View style={styles.container}>
        <Text style={styles.greeting}>Hello, {name}!</Text>
        <Text style={styles.subtitle}>You are logged in</Text>
        <TouchableOpacity style={styles.buttonOutline} onPress={handleLogout}>
          <Text style={styles.buttonOutlineText}>Log out</Text>
        </TouchableOpacity>
        <StatusBar style="auto" />
      </View>
    );
  }

  // Login form
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Foreningsliv</Text>

      <TextInput
        style={styles.input}
        placeholder="Email"
        placeholderTextColor="#999"
        value={email}
        onChangeText={setEmail}
        autoCapitalize="none"
        keyboardType="email-address"
      />

      <TextInput
        style={styles.input}
        placeholder="Password"
        placeholderTextColor="#999"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
      />

      {error && <Text style={styles.error}>{error}</Text>}

      <TouchableOpacity
        style={styles.button}
        onPress={handleLogin}
        disabled={loading}
      >
        {loading ? (
          <ActivityIndicator color="#fff" />
        ) : (
          <Text style={styles.buttonText}>Log in</Text>
        )}
      </TouchableOpacity>

      <StatusBar style="auto" />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#fff",
    alignItems: "center",
    justifyContent: "center",
    padding: 20,
  },
  title: {
    fontSize: 28,
    fontWeight: "bold",
    marginBottom: 32,
    color: "#333",
  },
  greeting: {
    fontSize: 28,
    fontWeight: "bold",
    color: "#333",
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 16,
    color: "#666",
    marginBottom: 24,
  },
  input: {
    width: "100%",
    maxWidth: 320,
    height: 48,
    borderWidth: 1,
    borderColor: "#ddd",
    borderRadius: 8,
    paddingHorizontal: 16,
    marginBottom: 12,
    fontSize: 16,
    backgroundColor: "#fafafa",
  },
  button: {
    width: "100%",
    maxWidth: 320,
    height: 48,
    backgroundColor: "#0066cc",
    borderRadius: 8,
    alignItems: "center",
    justifyContent: "center",
    marginTop: 8,
  },
  buttonText: {
    color: "#fff",
    fontSize: 16,
    fontWeight: "600",
  },
  buttonOutline: {
    width: "100%",
    maxWidth: 320,
    height: 48,
    borderWidth: 1,
    borderColor: "#cc0000",
    borderRadius: 8,
    alignItems: "center",
    justifyContent: "center",
  },
  buttonOutlineText: {
    color: "#cc0000",
    fontSize: 16,
    fontWeight: "600",
  },
  error: {
    color: "red",
    marginBottom: 8,
    fontSize: 14,
  },
});
