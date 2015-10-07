//#include "display.h"
#include <QDebug>

#include "dbus.h"

void hehe()
{
    qDebug() << "screenWidthChanged...";
}

inline QDebug operator<<(QDebug d, const QDBusObjectPath& p)
{
    QList<QVariant> argumentList;
    argumentList << QVariant::fromValue(p);
    d << p.path();
    return d;
}

using namespace dbus::common;
using namespace dbus::objects;

extern void body();
int main(int argc, char *argv[])
{
    QCoreApplication app(argc, argv);
    body();
    app.exec();
}
