import 'package:flutter/material.dart';

void main() => runApp(const WapasnamaApp());

/// Wapasnama — supplier claim register. Log every discrepancy raised against a
/// supplier and track credit-note promised vs received. Mirrors the Go journal
/// service (POST /log, GET /summary).
class WapasnamaApp extends StatelessWidget {
  const WapasnamaApp({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp(
        title: 'Wapasnama',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(colorSchemeSeed: const Color(0xFFB23A48), useMaterial3: true),
        home: const HomePage(),
      );
}

class Claim {
  final String supplier, item, kind;
  final double amount;
  bool received;
  Claim(this.supplier, this.item, this.kind, this.amount, {this.received = false});
}

/// openTotal sums the amounts of claims still awaiting a credit note.
double openTotal(List<Claim> claims) =>
    claims.where((c) => !c.received).fold(0.0, (s, c) => s + c.amount);

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _claims = <Claim>[];
  final _supplier = TextEditingController();
  final _item = TextEditingController();
  final _amount = TextEditingController();
  String _kind = 'short';

  void _add() {
    final amt = double.tryParse(_amount.text.trim()) ?? 0;
    if (_supplier.text.trim().isEmpty || amt <= 0) return;
    setState(() {
      _claims.insert(0, Claim(_supplier.text.trim(),
          _item.text.trim().isEmpty ? '—' : _item.text.trim(), _kind, amt));
      _supplier.clear();
      _item.clear();
      _amount.clear();
    });
  }

  @override
  Widget build(BuildContext context) {
    final open = openTotal(_claims);
    return Scaffold(
      appBar: AppBar(
        title: const Text('Wapasnama · supplier claims'),
        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
      ),
      body: Column(children: [
        Container(
          width: double.infinity,
          color: Theme.of(context).colorScheme.primaryContainer,
          padding: const EdgeInsets.all(16),
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            const Text('Open claims (awaiting credit note)'),
            Text('₹${open.toStringAsFixed(2)}',
                style: const TextStyle(fontSize: 30, fontWeight: FontWeight.bold)),
          ]),
        ),
        Padding(
          padding: const EdgeInsets.all(12),
          child: Row(children: [
            Expanded(child: TextField(controller: _supplier, decoration: const InputDecoration(labelText: 'Supplier', border: OutlineInputBorder()))),
            const SizedBox(width: 8),
            SizedBox(width: 100, child: TextField(controller: _amount, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: '₹', border: OutlineInputBorder()))),
          ]),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12),
          child: Row(children: [
            Expanded(child: TextField(controller: _item, decoration: const InputDecoration(labelText: 'Item', border: OutlineInputBorder()))),
            const SizedBox(width: 8),
            DropdownButton<String>(
              value: _kind,
              items: const [
                DropdownMenuItem(value: 'short', child: Text('short')),
                DropdownMenuItem(value: 'damaged', child: Text('damaged')),
                DropdownMenuItem(value: 'wrong-rate', child: Text('wrong-rate')),
                DropdownMenuItem(value: 'missing-scheme', child: Text('missing-scheme')),
              ],
              onChanged: (v) => setState(() => _kind = v ?? 'short'),
            ),
            const SizedBox(width: 8),
            FilledButton(onPressed: _add, child: const Text('Log')),
          ]),
        ),
        const Divider(),
        Expanded(
          child: ListView.builder(
            itemCount: _claims.length,
            itemBuilder: (_, i) {
              final c = _claims[i];
              return ListTile(
                title: Text('${c.supplier} · ${c.item}'),
                subtitle: Text('${c.kind} · ₹${c.amount.toStringAsFixed(2)}'),
                trailing: TextButton(
                  onPressed: () => setState(() => c.received = !c.received),
                  child: Text(c.received ? 'received ✓' : 'mark received'),
                ),
                leading: Icon(c.received ? Icons.check_circle : Icons.hourglass_empty,
                    color: c.received ? Colors.green : null),
              );
            },
          ),
        ),
      ]),
    );
  }
}
