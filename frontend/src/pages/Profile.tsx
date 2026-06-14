import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

interface User {
  username: string;
  points: number;
}

export default function Profile() {
  const { username } = useParams();

  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    fetch(`/users/${username}`)
      .then((res) => res.json())
      .then((data) => setUser(data));
  }, [username]);

  if (!user) {
    return <h1>Loading...</h1>;
  }

  return (
    <div>
      <h1>{user.username}</h1>
      <p>{user.points} points</p>
    </div>
  );
}