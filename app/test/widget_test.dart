import 'package:flutter_test/flutter_test.dart';

import 'package:wapasnama_app/main.dart';

void main() {
  test('openTotal sums only unreceived claims', () {
    final claims = [
      Claim('A', 'x', 'short', 300),
      Claim('B', 'y', 'damaged', 120, received: true),
      Claim('C', 'z', 'wrong-rate', 80),
    ];
    expect(openTotal(claims), 380);
  });

  testWidgets('renders the open-claims header', (tester) async {
    await tester.pumpWidget(const WapasnamaApp());
    expect(find.textContaining('Open claims'), findsOneWidget);
  });
}
