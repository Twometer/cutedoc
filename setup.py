from setuptools import setup, find_packages

VERSION = '1.0'

setup(
    name="cutedoc",
    version=VERSION,
    url='https://github.com/Twometer/cutedoc',
    license='MIT',
    description='Clean GitBook-inspired theme',
    author='Twometer',
    author_email='twometer@outlook.de',
    packages=find_packages(),
    include_package_data=True,
    entry_points={
        'mkdocs.themes': [
            'cutedoc = cutedoc_theme',
        ]
    },
    zip_safe=False
)
