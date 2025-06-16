const { useState } = React;

function App() {
  const [entries, setEntries] = useState([]);
  const [form, setForm] = useState({ id: null, name: '', phone: '', email: '' });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((f) => ({ ...f, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!form.name || !form.phone || !form.email) return;
    if (form.id === null) {
      const newEntry = { ...form, id: Date.now() };
      setEntries((prev) => [...prev, newEntry]);
    } else {
      setEntries((prev) => prev.map((e) => (e.id === form.id ? form : e)));
    }
    setForm({ id: null, name: '', phone: '', email: '' });
  };

  const handleEdit = (entry) => {
    setForm(entry);
  };

  const handleDelete = (id) => {
    setEntries((prev) => prev.filter((e) => e.id !== id));
    if (form.id === id) {
      setForm({ id: null, name: '', phone: '', email: '' });
    }
  };

  return (
    <div>
      <h1>Agenda Telefónica</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          name="name"
          placeholder="Nombre"
          value={form.name}
          onChange={handleChange}
        />
        <input
          type="text"
          name="phone"
          placeholder="Teléfono"
          value={form.phone}
          onChange={handleChange}
        />
        <input
          type="email"
          name="email"
          placeholder="Correo"
          value={form.email}
          onChange={handleChange}
        />
        <button type="submit">{form.id === null ? 'Agregar' : 'Actualizar'}</button>
      </form>

      <table>
        <thead>
          <tr>
            <th>Nombre</th>
            <th>Teléfono</th>
            <th>Correo</th>
            <th>Acciones</th>
          </tr>
        </thead>
        <tbody>
          {entries.map((entry) => (
            <tr key={entry.id}>
              <td>{entry.name}</td>
              <td>{entry.phone}</td>
              <td>{entry.email}</td>
              <td>
                <button onClick={() => handleEdit(entry)}>Editar</button>
                <button onClick={() => handleDelete(entry.id)}>Eliminar</button>
              </td>
            </tr>
          ))}
          {entries.length === 0 && (
            <tr>
              <td colSpan="4" style={{ textAlign: 'center' }}>
                Sin datos
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

ReactDOM.render(<App />, document.getElementById('root'));
